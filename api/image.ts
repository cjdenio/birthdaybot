import puppeteer from "puppeteer-core";
import nunjucks from "nunjucks";

import chromium from "chrome-aws-lambda";

import * as path from "path";

import { NowRequest, NowResponse } from "@vercel/node";

export default async (req: NowRequest, res: NowResponse) => {
  const browser = await puppeteer.launch({
    executablePath: await chromium.executablePath,
    args: chromium.args,
    headless: chromium.headless,
  });
  const page = await browser.newPage();
  await page.setViewport({
    width: 1920,
    height: 1080,
  });
  await page.setContent(
    nunjucks.render(path.join(__dirname, "_lib", "image.html"), {
      text: req.query.text,
      image: req.query.image,
      date: req.query.date,
    })
  );
  const screenshot = await page.screenshot();

  res.setHeader("Content-Type", "image/png");
  res.setHeader('Cache-Control', `public, immutable, no-transform, s-maxage=31536000, max-age=31536000`);
  
  res.send(screenshot);

  await browser.close();
};
