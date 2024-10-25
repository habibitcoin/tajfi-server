Hackathon Submission - Tajfi - A Nostr based Taproot Assets wallet powered by a Pocket Universe
<!-- All software-based projects submitted must be open source and freely available for public use -->

![ZBD x PlebLab Hackathon Image](https://pbs.twimg.com/media/GW2IHa2WYAE71ca?format=jpg&name=large)

# 🚀 TABConf 6 Hackathon Submission

## Project Name: Tajfi - A Nostr based Taproot Assets wallet powered by a Pocket Universe

### Team Members 👥

- Jad Mubaslat/Habibitcoin - Backend Developer
- Ruben Navarro - Full Stack Developer

### Project Description 📝

Tajfi is a Taproot Assets wallet that empowers users to manage virtual UTXOs (vUTXOs) while securely anchoring them to Bitcoin UTXOs. Leveraging Nostr for authentication and the Taproot Assets protocol, users maintain full custody over their assets with enhanced privacy. Marketplace functionality can also be built on top of this infrastructure, enabling users to trade PSBTs and manage Taproot Asset channels. PSBT marketplaces have already previously been implemented over Nostr (see [deezy.place](https://deezy.place)) and can be integrated into Tajfi to provide a seamless user experience, since Taproot Assets recycles the same PSBT tools as native Bitcoin transactions.

### Technical Implementation 💻

-   **Language:** Go (1.22+) & Next.js (React)
-   **Taproot Daemon Fork:** Custom fork ([habibitcoin/taproot-assets](https://github.com/habibitcoin/taproot-assets/tree/tajfi-fork)) with signature communication between `tapd` and `tajfi-server`. We had to fork the Taproot Assets Daemon to allow for the communication of signatures between the `tapd` and `tajfi-server`.
-   **Nostr Authentication:** Users log in and sign Taproot Assets transactions via Nostr pubkeys. The Tajfi universe operator is cryptographically **unable** to spend user funds.
-   **Demo Mode:** Preconfigured to fund invoices of amount 10 for LUSD on the [demo site](https://demo.tajfi.com), which speaks to the [demo api](https://api.tajfi.com).


### User Experience & Design 🎨

Tajfi offers seamless login with Nostr, transparent asset custody, and reduced reliance on centralized entities. It simplifies complex asset management, with backend control hidden from the user while ensuring complete asset ownership. User's can authenticate with any Nostr-compatible wallet that they have today. As long as the wallet can provide Schnorr signatures, it can be used to sign Taproot Assets transactions on Tajfi!


### Innovation & Creativity 💡

-   **Pocket Universe:** Creates isolated environments for users to manage and transact with Taproot Assets, without needing to worry about proof storage or revealing their onchain history to a blockchain, nor worry about managing onchain Bitcoin UTXOs.
-   **Privacy-focused custody:** Anchors vUTXOs to Bitcoin UTXOs while giving the user control of assets. Only the service provider has access to the underlying Bitcoin UTXOs, not the user's Taproot Assets.
-   **Planned 2-of-2 MuSig2 UTXO Control:** Ensures that the Tajfi server cannot accidentally burn any Taproot Assets, providing a more secure environment for users.
-   **Marketplace Potential:** Enables the creation of PSBT marketplaces and Taproot Asset channel management services, opening up new opportunities for value-added services. By leveraging Nostr to relay PSBT offers ([as described by OpenOrdex](https://github.com/orenyomtov/openordex/blob/44581ec727c439c15178413b1d46c8f6176f253a/NIP.md?plain=1#L2) using `kind 802`), Taproot Assets can be exchanged for on-chain BTC in a censorship-resistant manner. The skeleton for this is already in place on the [demo site's marketplace](https://demo.tajfi.com/wallet/marketplace).
-  **Scalability Enhancements:** By aggregating proofs of vUTXOs, the onchain footprint of Taproot Assets transactions can be reduced, making it more scalable and cost efficient. A single Bitcoin UTXO can represent multiple Taproot Asset users' state transitions.


### Potential Impact 🌍

Tajfi provides a decentralized alternative to centralized wallets and ERC-20 tokens, promoting privacy and asset control. It opens new opportunities for services like Taproot Asset channel management and PSBT marketplaces. This enables Bitcoin to continue to connect a broader financial ecosystem by operating as a settlement and liquidity layer. The more interconnected the Bitcoin network becomes, the more valuable it is to the world.

We hope that this technology can enable people who can't afford to have their assets seized by foreign governments, nor divulge their privacy, to still be able to secure their wealth and transact with the world and the growing Bitcoin economy.


### Business Model 💼

A Tajfi pocket universe operator has a few different ways to monetize their service:
- Fee collection on Taproot Asset transactions (e.g., 0.1% of the transaction amount)
- Charging for additional services, such as Taproot Asset channel management or PSBT marketplace services
- Charging token issuers for native application support (but this feels like shitcoining, so we're not too excited about this one)
- Advertising and partnerships with other Bitcoin services or related financial services


### Demo Video 🎥

[Link to your 1-minute demo video]


### GitHub Repository 📂

-   [Tajfi Backend](https://github.com/habibitcoin/taproot-assets/tree/tajfi-fork)
-   [Tajfi Frontend](https://github.com/topether21/tajfi-web)
-   [Taproot Assets Custom Fork](https://github.com/habibitcoin/taproot-assets/tree/tajfi-fork)

### Additional Resources 📚

-   [Demo Site](https://demo.tajfi.com)
-   [Demo API](https://api.tajfi.com)
-   [Taproot Assets FAQ](https://docs.lightning.engineering/the-lightning-network/taproot-assets/faq)
-   [Taproot Assets Trustless Swap](https://docs.lightning.engineering/the-lightning-network/taproot-assets/trustless-swap)


### Future Plans 🔮

Post-hackathon, the project aims to implement MuSig2 for secure the anchor Bitcoin UTXOs, ensuring that no single party can accidentally burn assets. We then plan to follow up with the implementation of the PSBT marketplace. Afterwards, we will explore additional improvements, such as Lightning Network support/Taproot Asset channels, and/or proof aggregation of vUTXOs to reduce the onchain footprint of Taproot Assets transactions.

### Feedback for Organizers 📣

Nothing I can think of really! First time at TABConf and has to be one of my favorites. Signal:noise ratio is great.
