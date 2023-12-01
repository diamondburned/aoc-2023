{ pkgs ? import <nixpkgs> {} }:

let inputpaste = pkgs.writeShellScriptBin "inputpaste" ''
	set -e
	wl-paste > input
	echo Pasted to ./input:
	head -n10 input
'';

in pkgs.mkShell {
	buildInputs = with pkgs; [
		go
		gopls
		gotools
		zls
		zig
		inputpaste
	];
}
