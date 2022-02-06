import flag from './flag.txt'

function sleep(time) {
  return new Promise(resolve => {
    setTimeout(resolve, time)
  })
}

export default {
  id: 'shadow',
  name: 'shadow',
  urlRegex: /^https:\/\/shadow\.mc\.ax\//,
  timeout: 10000,
  handler: async (url, ctx) => {
    const page = await ctx.newPage()
    await page.goto('https://shadow.mc.ax/', { timeout: 3000, waitUntil: 'domcontentloaded' })
    await page.evaluate(secret => localStorage.setItem('secret', secret), flag)
    await page.goto(url, { timeout: 3000, waitUntil: 'domcontentloaded' })
    await sleep(3000)
  }
}
