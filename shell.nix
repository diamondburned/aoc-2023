let ov = self: super: {
	go = super.go_1_19;
};

in { pkgs ? import <nixpkgs> { overlays = [ ov ]; } }:

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
