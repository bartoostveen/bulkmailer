{
  buildGoModule,
  lib,
  ...
}:

buildGoModule (finalAttrs: {
  pname = "bulkmailer";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-rsNaS0RSEaVDcb7o8Bd0s+1jbGEb8SLNvz1+orDf/RI=";

  meta = {
    description = "A very primitive template renderer and sender, written in Go.";
    homepage = "https://git.bartoostveen.nl/bart/bulkmailer";
    license = lib.licenses.gpl3Only;
    mainProgram = "bulkmailer";
  };
})
