Now I need to work on creating PGP Key for user and assign it in database

TODO:

1. PGP Form validation when user, enters the PGP key and when storing in the database
   (How should I protect the input, how should i store the PGP key into the database what library should I use and etc)

2. CHATTING FUNCTIONALITY after complete PGP functionality (between users)

DOCKER: I can try to host the application on onion and try to test the application

Example of chat component: <https://v0.dev/chat/decentralized-chat-platform-M8TWUSTLj6A>

⚠️ Important:
You are the only one with access to your private key.
If you lose your key or forget your passphrase, your messages cannot be recovered.
This is a privacy feature — not even we can help you regain access

<http://securityheaders.com/>
<https://developer.mozilla.org/en-US/observatory>

Implementing encryption.

Private key should live encrypted in session key
Passphrase for encrypted key also in session key
