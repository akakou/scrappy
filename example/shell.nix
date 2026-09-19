{ pkgs ? import ../nix/pkgs.nix }:

import ../shell.nix { inherit pkgs; }
