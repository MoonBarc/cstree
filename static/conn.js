const statuses = {
    "connecting...": "#8f8f8f",
    "connected": "#75ff7c",
    "paused": "#ffd342",
    "down": "#ff5757"
}

const TOTAL_LIGHTS = 150

function setStatus(to) {
    const s = document.getElementById("status")
    const d = document.getElementById("statusdot")
    s.textContent = to
    d.style.backgroundColor = statuses[to]
}

/**
 * @param {HTMLImageElement} trimage
 * @param {CanvasRenderingContext2D} ctx
 * @param {string[]} lights
 */
function render(trimage, ctx, lights) {
    const lrad = 8;
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height)
    ctx.drawImage(trimage, 20, 0)

    let rowi = 0;
    let row = 1;
    for (let light of lights) {
        if ((rowi + 1) % row == 0) {
            row ++;
            rowi = 0;
        }
        ctx.beginPath()
        ctx.fillStyle = light
        ctx.ellipse(
            (trimage.width / 2) + 37 - (row * lrad) + (rowi * lrad * 2),
            (row * lrad * 2.8) - ((row-rowi-1) * 1.5),
            lrad, lrad, 0, 0, 360)
        ctx.fill()
        rowi++
    }
}

async function main() {
    /**
     * @type {HTMLCanvasElement}
     */
    const canvas = document.getElementById("treestatus")
    const ctx = canvas.getContext("2d")
    const trimage = new Image()
    let loaded = false
    const events = new EventSource("/events")

    events.onmessage = (e) => {
        if (!loaded) {
            return
        }
        const d = JSON.parse(e.data)
        if (d.Ty == "render") {
            const colors = d.Data.colors
            const bytes = Uint8Array.from(atob(colors), c => c.charCodeAt(0))
            const stringColors = []
            for (let i = 0; i < bytes.length; i += 4) {
                let str = "#"
                for (let j = 1; j < 4; j++) {
                    str += bytes[i + (3-j)].toString(16).padStart(2, "0")
                }
                stringColors[149-(i/4)] = str
            }
            render(trimage, ctx, stringColors)
        }
        if (d.Ty == "info") {
            document.getElementById("codetitle").innerText = `"${d.Data.title}"`
            document.getElementById("author").innerText = d.Data.author
        }
    }

    events.onopen = (_e) => {
        setStatus("connected")
    }

    events.onerror = (_e) => {
        setStatus("connecting...")
    }

    trimage.src = "/tree.png"
    trimage.addEventListener("load", () => {
        canvas.height = trimage.height
        canvas.width = trimage.width + 40
        const lightstates = []
        const colors = ["#ff00ff", "#00ff00", "#ff0000", "#0f0fff"]
        for (let i = 0; i < TOTAL_LIGHTS; i++) {
            lightstates[i] = colors[Math.floor(Math.random() * colors.length)]
        }
        render(trimage, ctx, lightstates)
        loaded = true
    })
}

document.addEventListener("DOMContentLoaded", main);