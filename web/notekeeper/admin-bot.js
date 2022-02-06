export default {
  id: 'notekeeper',
  name: 'notekeeper',
  timeout: 25000,
  handler: async (url, ctx) => {
    const page = await ctx.newPage();

    await page.goto(url);
    await page.waitForTimeout(4000);

    // JWT for user "admin", signed with key from challenge.yaml
    const jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIn0.coi9nd8-GchMdBOnKtyE5cp4Ms4zLyjowB4-Y_0BnvU";
    await page.setCookie({
      name: 'session',
      value: jwt,
      domain: "notekeeper.mc.ax",
      httpOnly: true
    });
    await page.goto("https://notekeeper.mc.ax/home");
    await page.waitForTimeout(4000);

    await page.evaluate(() => {
      document.querySelector("#logout") && document.querySelector("#logout").click();
    });
    await page.waitForNavigation();
  },
};
