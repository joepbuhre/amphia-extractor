# { pkgs ? import <nixpkgs> {} }:
# let
# in
#   pkgs.mkShell {
#     buildInputs = [
#       pkgs.nodejs_20
#       pkgs.go
#       pkgs.air
#       pkgs.go-swag
#       pkgs.gopls
#     ];

#     # export GOPATH="/home/jbuhre/development/joepbuhre/snappic/backend/.gopath"
#     # export PATH="$PATH:/home/jbuhre/development/joepbuhre/snappic/backend/.gopath/bin"
#   shellHook = ''
#     export GOROOT="${pkgs.go}"
#   '';
# }

{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  PROJECT_ROOT = builtins.toString ./.;

  buildInputs = [
    nodejs_20
    go
    gotools
    gopls
    go-outline
    gopkgs
    gocode-gomod
    godef
    golint
    air
    delve
    http-server
    sqlc
    sqlite
  ];
  shellHook = ''
    export CUR_DIR=/home/jbuhre/development/joepbuhre/money-manager
    export PATH="$PATH:$CUR_DIR/backend/.gopath/bin"
    export GOPATH="$CUR_DIR/backend/.gopath"

    start() {
      cd "$CUR_DIR"
      # Check if at least one argument is provided
      if [ $# -lt 1 ]; then
          echo "Usage: start (frontend | backend)" 
          return 1
      fi

      # Get the argument
      service=$1

      # Define the function to start the service
      start_service() {
          local dir=$1
          if [ "$(basename "$PWD")" != "$dir" ]; then
              echo "Changing directory to $dir"
              cd $dir || { echo "Failed to change directory to $dir"; return 1; }
          fi
          echo "Starting $dir..."
      }

      # Start the frontend or backend service based on the argument
      case $service in
          "frontend")
              start_service "frontend"
              npm run dev
              ;;
          "backend")
              start_service "backend"
              air
              ;;
          *)
              echo "Invalid argument. Usage: start (frontend | backend)"
              return 1
              ;;
      esac
  }

  create-migration() {
    result=$(goose --dir $CUR_DIR/backend/sql/migrations create $1 sql 2>&1)
    path=$(echo $result | sed 's/.*Created new file: //')
    code "$path" 
  }
    alias start-backend="cd backend && air"
    alias start-frontend="cd frontend && npm run dev"
    code .

  '';

}
