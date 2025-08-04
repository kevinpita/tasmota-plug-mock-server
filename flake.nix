{
  description = "Go development environment";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };

      gitignoreContent = ''
        .direnv

        .vscode/
        .idea/
      '';
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          git

          go
          gofumpt
          golangci-lint
          gopls

          pre-commit
        ];
        shellHook = ''
          projectName=$(basename "$(pwd)" | sed 's/ /-/g')
          if [ ! -f go.mod ]; then
            go mod init "github.com/kevinpita/$projectName"
          fi

          if [ ! -d .git ]; then
            git init
            echo "${gitignoreContent}" > .gitignore
            git add .
            git commit -m "Initial commit"
            pre-commit install --install-hooks
            git remote add origin "git@github.com:kevinpita/$projectName.git"
          fi
        '';
      };
    };
}
