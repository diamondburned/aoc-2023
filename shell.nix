let ov = self: super: {
	go = super.go_1_19;
};

in { pkgs ? import <nixpkgs> { overlays = [ ov ]; } }:

pkgs.mkShell {
	buildInputs = with pkgs; [
		go
		gopls
		gotools
		zls
		zig
	];
}
