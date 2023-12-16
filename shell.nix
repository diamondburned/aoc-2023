{ pkgs ? import <nixpkgs> {} }:

let
	inputpaste = pkgs.writeShellScriptBin "inputpaste" ''
		set -e
		wl-paste > input
		echo Pasted to ./input:
		head -n10 input
	'';

	overrideGo = go:
		go.overrideAttrs (old: rec {
			pname = old.pname + "-aoc";
			version = "a94347a";

			src = pkgs.fetchgit {
				url = "https://go.googlesource.com/go";
				rev = version;
				hash = "sha256-7UmPaQKWYiilLoSpkgsMUjCnQVzF+IHe6zzEUwjsugs=";
				leaveDotGit = true;
			};
		
			nativeBuildInputs = (old.nativeBuildInputs or []) ++ (with pkgs; [
				git # so Go can find its version
			]);

			patches = with pkgs; [
				(substituteAll {
					src = <nixpkgs/pkgs/development/compilers/go/iana-etc-1.17.patch>;
					iana = iana-etc;
				})
				# Patch the mimetype database location which is missing on NixOS.
				# but also allow static binaries built with NixOS to run outside nix
				(substituteAll {
					src = <nixpkgs/pkgs/development/compilers/go/mailcap-1.17.patch>;
					inherit mailcap;
				})
				# prepend the nix path to the zoneinfo files but also leave the original value for static binaries
				# that run outside a nix server
				(substituteAll {
					src = <nixpkgs/pkgs/development/compilers/go/tzdata-1.19.patch>;
					inherit tzdata;
				})
				(<nixpkgs/pkgs/development/compilers/go/remove-tools-1.11.patch>)
				(pkgs.fetchurl {
					url = "https://gist.githubusercontent.com/diamondburned/7002c13a4417f083c905e5a67061c9a5/raw/952c69bd92e32d78a8c26ebfe6ea5d3c96e3515a/go_no_vendor_checks-1.22.patch";
					hash = "sha256-u33SHl87PE7ZX/2aZKzV8sOI17RqrzcmIWlMv81z3Bk=";
				})
				(pkgs.fetchpatch {
					url = "https://go-review.googlesource.com/changes/go~510541/revisions/18/patch?zip";
					hash = "sha256-pvTQPoGay99fySyKWIicHTuOLE+uGQ0m1/50w59LJsc=";
					decode = "${pkgs.gzip}/bin/gzip -d";
				})
			];
		});

	patchedPkgs = pkgs.extend (self: super: {
		go = overrideGo super.go;
		buildGoModule = super.buildGoModule.override {
			inherit (self) go;
		};
		buildGo121Module = super.buildGo121Module.override {
			inherit (self) go;
		};
	});
in

pkgs.mkShell {
	buildInputs = with patchedPkgs; [
		go
		gopls
		gotools
		go-tools
	];

	GOEXPERIMENT = "loopvar,range";
}
