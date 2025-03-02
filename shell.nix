{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
	buildInputs = [
        pkgs.sqlite
        pkgs.ffmpeg-full
        pkgs.go-swag
        pkgs.go
	];
}

