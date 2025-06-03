import {
  loadPrivateKey,
  decryptMessage,
  clearSession,
  getPrivateKeyFromSession,
} from "./usePgp.js";

window.handleKeyUpload = async function () {
  const file = document.getElementById("privateKeyFile").files[0];
  if (!file) {
    alert("Please select a private key file.");
    return;
  }

  try {
    await loadPrivateKey(file);
    updateKeyUI(true);
  } catch (err) {
    console.error(err);
    alert("❌ Failed to load private key.");
  }
};

window.handleClearKey = function () {
  clearSession();
  updateKeyUI(false);
};

function updateKeyUI(isLoaded) {
  const uploadArea = document.getElementById("key-upload-area");
  const clearArea = document.getElementById("key-clear-area");
  const status = document.getElementById("key-status");
  const expiryDisplay = document.getElementById("key-expiry");

  if (isLoaded) {
    uploadArea.style.display = "none";
    clearArea.style.display = "block";
    status.innerText = "✅ Private key is loaded and active.";
    startCountdownTimer();
  } else {
    uploadArea.style.display = "block";
    clearArea.style.display = "none";
    status.innerText = "⚠️ No private key loaded.";
    expiryDisplay.innerText = "";
    if (window._countdownInterval) {
      clearInterval(window._countdownInterval);
    }
  }
}

function startCountdownTimer() {
  const expiryDisplay = document.getElementById("key-expiry");

  function updateCountdown() {
    const expiry = parseInt(sessionStorage.getItem("pgp_expiry") || "0", 10);
    const now = Date.now();
    const secondsLeft = Math.floor((expiry - now) / 1000);

    if (isNaN(expiry) || secondsLeft <= 0) {
      clearSession();
      updateKeyUI(false);
      clearInterval(window._countdownInterval);
      return;
    }

    expiryDisplay.innerText = `⏳ Expires in ${secondsLeft} seconds.`;
  }

  updateCountdown();
  clearInterval(window._countdownInterval);
  window._countdownInterval = setInterval(updateCountdown, 1000);
}

window.addEventListener("DOMContentLoaded", () => {
  const isLoaded = getPrivateKeyFromSession();
  updateKeyUI(isLoaded);
});
