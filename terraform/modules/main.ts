import { App} from "cdktf";
import { AwsProvider } from "@cdktf/provider-aws/lib/provider";
import {RandomProvider} from "@cdktf/provider-random/lib/provider"
import { HappyServiceModule } from "./service/main";

const app = new App();
const stack = new HappyServiceModule(app, "happy-dns");
new AwsProvider(stack, "AWS", {
  region: "us-west-2",
});
new RandomProvider(stack, "Random", {
})
app.synth();
