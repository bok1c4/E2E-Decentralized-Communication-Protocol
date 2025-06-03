import { getDecryptedPrivateKey } from "./usePgp.js";

async function decryptAllMessages() {
  const messages = document.querySelectorAll(".pgp-message");
  for (const el of messages) {
    const armored = el.innerText.trim();
    try {
      const msg = await openpgp.readMessage({ armoredMessage: armored });
      const privateKey = await getDecryptedPrivateKey();
      const { data: plaintext } = await openpgp.decrypt({
        message: msg,
        decryptionKeys: privateKey,
      });
      el.innerText = plaintext;
    } catch (err) {
      console.warn("Decryption failed:", err);
      el.innerText = "[🔒 Unable to decrypt]";
    }
  }
}

document.addEventListener("htmx:afterSwap", (e) => {
  if (e.target.id === "messages") {
    decryptAllMessages();
  }
});
