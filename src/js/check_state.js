window.onload = async function updateStates() {
    let completedStates = 0
    if ("session" in sessionStorage) {
        let resp = await fetch(`/api/${sessionStorage.session}`);
        let data = await resp.json();
        if (data.discord_name !== "") {
            document.getElementById("discord-text").innerText = "You are logged in as " + data.discord_name
            completedStates += 1
        }
        if (resp.status === 511) {
            document.getElementById("info-message").innerHTML =
                "You don't have a ETF2L profile, please visit <a href='https://etf2l.org'>ETF2L</a> and register"
            return
        }
        if (data.steam_name !== "") {
            document.getElementById("steam-text").innerHTML = "You are logged in as " + data.steam_name
            completedStates += 1
        }
        if (completedStates === 2) {
            document.getElementById("info-message").innerText =
                "Thank you, now you have full access to ETF2L discord server!"
        }
    }
}

document.getElementById("discord-button").onclick = async function authorise_discord() {
    if (sessionStorage.session === undefined) {
        sessionStorage.session = await register()
    }
    window.location.replace(`/auth/discord/?state=${sessionStorage.session}`);
}

document.getElementById("steam-button").onclick = async function authorise_steam() {
    if (sessionStorage.session === undefined) {
        sessionStorage.session = await register()
    }
    window.location.replace(`/auth/steam/?state=${sessionStorage.session}`);
}

async function register() {
    let response = await fetch("/api/register", {
        method: "post"
    });
    let data = await response.json();
    return data.state
}

async function checkEtf2lProfile(steam_id) {
    let resp = await fetch(`https://api.etf2l.org/player/${steam_id}.json`)
    if (resp.status === 500) {
        return "unknown"
    }
    let data = await resp.json()
    return data.name
}