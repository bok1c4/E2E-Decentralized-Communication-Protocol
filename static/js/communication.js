import { getDecryptedPrivateKey } from "./usePgp.js";

window.addEventListener("DOMContentLoaded", () => {
  const form = document.querySelector("form[action^='/channel/']");
  if (!form) return;

  form.addEventListener("submit", async (e) => {
    const meta = document.getElementById("pgp-channel-meta");
    const textarea = form.querySelector("textarea[name='content']");
    const message = textarea.value.trim();
    if (!message) return;

    if (meta && meta.dataset.isDirect === "true") {
      e.preventDefault();

      const publicKeyArmored = meta.dataset.recipientPublicKey;
      try {
        const pubKey = await openpgp.readKey({ armoredKey: publicKeyArmored });

        const encrypted = await openpgp.encrypt({
          message: await openpgp.createMessage({ text: message }),
          encryptionKeys: pubKey,
        });

        textarea.value = encrypted;
        form.submit(); // submit after encryption
      } catch (err) {
        alert("❌ Failed to encrypt message.");
        console.error(err);
      }
    }
  });
});
