{
  buildGoModule,
  lib,
  ...
}:

buildGoModule (finalAttrs: {
  pname = "bulkmailer";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-+BVsB1jC1uNFj9qb0wI2Ifpy9Mva5a0iz6iRbjvOPa8=";

  meta = {
    description = "A very primitive template renderer and sender, written in Go.";
    homepage = "https://git.bartoostveen.nl/bart/bulkmailer";
    license = lib.licenses.gpl3Only;
    mainProgram = "bulkmailer";
  };
})
