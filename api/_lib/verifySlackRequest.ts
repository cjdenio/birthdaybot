import { NowRequest } from "@vercel/node";
import * as crypto from "crypto";

export default ({
  timestamp,
  signingSecret,
  signature,
  req,
}: {
  timestamp?: string;
  signingSecret?: string;
  signature?: string;
  req: NowRequest;
}) => {
  return new Promise((resolve, reject) => {
    let bodyArr: Buffer[] = [];
    req
      .on("data", (chunk) => {
        bodyArr.push(chunk);
      })
      .on("end", () => {
        const body = Buffer.concat(bodyArr).toString();

        const base = `v0:${timestamp}:${body}`;

        const hmac = crypto.createHmac("sha256", signingSecret);
        hmac.update(base);

        resolve(
          crypto.timingSafeEqual(
            Buffer.from("v0=" + hmac.digest("hex")),
            Buffer.from(signature)
          )
        );
      });
  });
};
