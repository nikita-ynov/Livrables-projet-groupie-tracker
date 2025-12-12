import { createWeb3Modal, defaultWagmiConfig } from 'https://esm.sh/@web3modal/wagmi'
import { reconnect, watchAccount, signMessage, getAccount, disconnect } from 'https://esm.sh/@wagmi/core'
import { mainnet, sepolia } from 'https://esm.sh/@wagmi/core/chains'

const projectId = '3a8170812b534d0ff9d794f19a901d64'

const metadata = {
  name: 'YNFT',
  description: 'YNFT Marketplace',
  url: 'http://localhost:8080', 
  icons: ['https://avatars.githubusercontent.com/u/37784886']
}

const chains = [mainnet, sepolia]
const config = defaultWagmiConfig({ chains, projectId, metadata })
const modal = createWeb3Modal({ wagmiConfig: config, projectId })

let isSigningIn = false;

// LOGIN SERVEUR
async function loginToServer(address) {
    if (isSigningIn) return;
    isSigningIn = true;
    try {
        const nonce = "Login à YNFT: " + Date.now();
        console.log("Demande de signature...");
        
        const signature = await signMessage(config, { message: nonce });

        const response = await fetch('/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ address, signature, nonce })
        });

        if (response.ok) {
            console.log("✅ Authentifié côté serveur");
            window.location.reload(); 
        } else {
            console.error("❌ Erreur authentification serveur");
            disconnect(config);
        }
    } catch (error) {
        console.error("Erreur login:", error);
        disconnect(config);
    } finally {
        isSigningIn = false;
    }
}

// SETUP BOUTONS CONNECT
function setupConnectButtons() {
    const buttons = document.querySelectorAll('.login__btn');
    buttons.forEach(btn => {
        const newBtn = btn.cloneNode(true);
        btn.parentNode.replaceChild(newBtn, btn);
        newBtn.addEventListener('click', () => { modal.open(); });
    });
}

document.addEventListener('DOMContentLoaded', () => {
    setupConnectButtons();
});

watchAccount(config, {
    onChange: async (account) => {
        const buttons = document.querySelectorAll('.login__btn');
        if (account.isConnected && account.address) {
            const addr = account.address;
            buttons.forEach(btn => {
                btn.innerText = `${addr.substring(0, 4)}...${addr.substring(addr.length - 4)}`;
                
                // --- MODIFICATION ICI ---
                // On ajoute la classe CSS pour gérer la couleur via le fichier styles.css
                btn.classList.add("connected");
                btn.style.backgroundColor = ""; // On nettoie les styles inline au cas où
            });
            
            const hasCookie = document.cookie.split(';').some((item) => item.trim().startsWith('session_token='));
            if (!hasCookie && !isSigningIn) {
                await loginToServer(addr);
            }
        } else {
            buttons.forEach(btn => {
                btn.innerText = "Connect";
                
                // --- MODIFICATION ICI ---
                // On retire la classe CSS quand déconnecté
                btn.classList.remove("connected");
                btn.style.backgroundColor = ""; 
            });
        }
    }
});

reconnect(config);