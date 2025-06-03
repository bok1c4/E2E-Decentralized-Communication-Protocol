import * as openpgp from "https://unpkg.com/openpgp@5.10.1/dist/openpgp.min.mjs";

const sessionKeys = {
  encryptedKey: "pgp_encrypted_key",
  aesKey: "pgp_aes_key",
  expiry: "pgp_expiry",
};

function arrayBufferToBase64(buffer) {
  return btoa(String.fromCharCode(...new Uint8Array(buffer)));
}

function base64ToArrayBuffer(base64) {
  const binary = atob(base64);
  const len = binary.length;
  const bytes = new Uint8Array(len);
  for (let i = 0; i < len; i++) bytes[i] = binary.charCodeAt(i);
  return bytes.buffer;
}

export async function loadPrivateKey(file, armoredKeyString) {
  const armoredKey = file ? await file.text() : armoredKeyString;
  const privateKey = await openpgp.readPrivateKey({ armoredKey });
  const rawKey = new TextEncoder().encode(privateKey.armor());

  const aesKey = crypto.getRandomValues(new Uint8Array(32));
  const cryptoKey = await crypto.subtle.importKey(
    "raw",
    aesKey,
    "AES-GCM",
    false,
    ["encrypt", "decrypt"],
  );
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const ciphertext = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv },
    cryptoKey,
    rawKey,
  );

  sessionStorage.setItem(
    sessionKeys.encryptedKey,
    JSON.stringify({
      iv: arrayBufferToBase64(iv),
      ciphertext: arrayBufferToBase64(ciphertext),
    }),
  );
  sessionStorage.setItem(sessionKeys.aesKey, arrayBufferToBase64(aesKey));
  sessionStorage.setItem(sessionKeys.expiry, Date.now() + 10 * 60 * 1000); // 10 min expiry
}

export async function getDecryptedPrivateKey() {
  const now = Date.now();
  const expiry = parseInt(
    sessionStorage.getItem(sessionKeys.expiry) || "0",
    10,
  );
  if (now > expiry) throw new Error("Private key session expired.");

  const { ciphertext, iv } = JSON.parse(
    sessionStorage.getItem(sessionKeys.encryptedKey),
  );
  const aesKey = base64ToArrayBuffer(
    sessionStorage.getItem(sessionKeys.aesKey),
  );
  const cryptoKey = await crypto.subtle.importKey(
    "raw",
    aesKey,
    "AES-GCM",
    false,
    ["decrypt"],
  );
  const decryptedBytes = await crypto.subtle.decrypt(
    { name: "AES-GCM", iv: base64ToArrayBuffer(iv) },
    cryptoKey,
    base64ToArrayBuffer(ciphertext),
  );
  const armored = new TextDecoder().decode(decryptedBytes);
  return openpgp.readPrivateKey({ armoredKey: armored });
}

export function getPrivateKeyFromSession() {
  try {
    const encrypted = sessionStorage.getItem(sessionKeys.encryptedKey);
    const aesKey = sessionStorage.getItem(sessionKeys.aesKey);
    const expiry = parseInt(
      sessionStorage.getItem(sessionKeys.expiry) || "0",
      10,
    );
    const now = Date.now();

    if (!encrypted || !aesKey || now > expiry) {
      return false;
    }

    return true;
  } catch (error) {
    console.error("Error checking private key session:", error);
    return false;
  }
}

export async function decryptMessage(armoredMessage) {
  const privateKey = await getDecryptedPrivateKey();
  const message = await openpgp.readMessage({ armoredMessage });
  const { data: plaintext } = await openpgp.decrypt({
    message,
    decryptionKeys: privateKey,
  });
  return plaintext;
}

export function clearSession() {
  sessionStorage.removeItem(sessionKeys.encryptedKey);
  sessionStorage.removeItem(sessionKeys.aesKey);
  sessionStorage.removeItem(sessionKeys.expiry);
}
