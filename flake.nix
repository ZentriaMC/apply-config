{
  description = "apply-config";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    let
      supportedSystems = [
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
      ];
    in
    flake-utils.lib.eachSystem supportedSystems (system:
      let
        rev = if (self ? rev) then self.rev else null;
        pkgs = nixpkgs.legacyPackages.${system};
      in
      rec {
        packages.apply-config = pkgs.callPackage ./default.nix { };

        devShell = pkgs.mkShell {
          nativeBuildInputs = [
            pkgs.go
            pkgs.golangci-lint
            pkgs.gopls
          ];
        };
      });
}
