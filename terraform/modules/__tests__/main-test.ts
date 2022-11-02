// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0
import "cdktf/lib/testing/adapters/jest"; // Load types for expect matchers
import { TerraformStack, Testing } from "cdktf";
import { HappyDNS } from "../main"
import { Route53Record } from "@cdktf/provider-aws/lib/route53-record";
describe("happy DNS", () => {
  it("create an AWS Route53 record", () => {
    const construct = Testing.synthScope((scope) => {
      new HappyDNS(scope, "my-dns", "albName", "albZoneId", "zone.czi.technology", "zoneId", "stackName", "appName");
    })
    expect(construct).toHaveResourceWithProperties(Route53Record, {
      type: "A",
      name: "stackName-appName.zone.czi.technology",
      zone_id: "zoneId", // weird, it changes from camelCase to snake_case otherwise the test will fail
      alias: [{
        name: "albName",
        zone_id: "albZoneId",
        evaluate_target_health: false,
      }]
    });
  });

  it("it looks like it did before", () => {
    const app = Testing.app();
    const stack = new TerraformStack(app, "test-happy-dns");
    new HappyDNS(stack, "my-dns", "albName", "albZoneId", "zone.czi.technology", "zoneId", "stackName", "appName");
    expect(Testing.synth(stack)).toMatchSnapshot();
  });


  //   it("Tests a combination of resources", () => {
  //     expect(
  //       Testing.synthScope((stack) => {
  //         new TestDataSource(stack, "test-data-source", {
  //           name: "foo",
  //         });

  //         new TestResource(stack, "test-resource", {
  //           name: "bar",
  //         });
  //       })
  //     ).toMatchInlineSnapshot();
  //   });
  // });

  // describe("Checking validity", () => {
  //   it("check if the produced terraform configuration is valid", () => {
  //     const app = Testing.app();
  //     const stack = new TerraformStack(app, "test");

  //     new TestDataSource(stack, "test-data-source", {
  //       name: "foo",
  //     });

  //     new TestResource(stack, "test-resource", {
  //       name: "bar",
  //     });
  //     expect(Testing.fullSynth(app)).toBeValidTerraform();
  //   });

  //   it("check if this can be planned", () => {
  //     const app = Testing.app();
  //     const stack = new TerraformStack(app, "test");

  //     new TestDataSource(stack, "test-data-source", {
  //       name: "foo",
  //     });

  //     new TestResource(stack, "test-resource", {
  //       name: "bar",
  //     });
  //     expect(Testing.fullSynth(app)).toPlanSuccessfully();
  //   });
  // });
});
