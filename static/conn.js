const statuses = {
    "connecting...": "#8f8f8f",
    "connected": "#75ff7c",
    "paused": "#ffd342",
    "down": "#ff5757"
}

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
    ctx.drawImage(trimage, 0, 0)

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
            (trimage.width / 2) - (row * lrad) + (rowi * lrad * 5),
            row * lrad * 4.1,
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
    trimage.src = "/tree.png"
    trimage.addEventListener("load", () => {
        canvas.height = trimage.height
        canvas.width = trimage.width
        render(trimage, ctx, ["#ff0000", "#ff00ff", "#00ff00", "#ff0000", "#ff0000", "#ff0000", "#ff0000", "#0f0fff"])
    })
}

document.addEventListener("DOMContentLoaded", main);