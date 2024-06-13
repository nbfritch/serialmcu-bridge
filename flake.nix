{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let pkgs = nixpkgs.legacyPackages.aarch64-linux.pkgs;
    in {
      devShells.aarch64-linux.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          python311
          python311Packages.pyserial
          bashInteractive
        ];
      };
    };
}
