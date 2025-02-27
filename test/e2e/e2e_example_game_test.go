package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

func TestE2eExampleGame(t *testing.T) {
	t.Parallel()

	// Our SSH username, will change based on which AMI we use
	username := "ubuntu"

	// Copy the terraform folder to a temp directory so we can run multiple tests in parallel
	tmpFolder := teststructure.CopyTerraformFolderToTemp(t, "..", "tf/public-ec2-instance")

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer teststructure.RunTestStage(t, "TEARDOWN", func() {
		teardown(t, tmpFolder)
	})

	// Deploy the terraform infra
	teststructure.RunTestStage(t, "SETUP", func() {
		setup(t, tmpFolder)
	})

	// Upload the Zarf artifacts
	teststructure.RunTestStage(t, "UPLOAD", func() {
		terraformOptions := teststructure.LoadTerraformOptions(t, tmpFolder)
		keyPair := teststructure.LoadEc2KeyPair(t, tmpFolder)

		syncFileToRemoteServer(t, terraformOptions, keyPair, username, "../../build/zarf", fmt.Sprintf("/home/%s/build/zarf", username), "0700")
		syncFileToRemoteServer(t, terraformOptions, keyPair, username, "../../build/zarf-init.tar.zst", fmt.Sprintf("/home/%s/build/zarf-init.tar.zst", username), "0600")
		syncFileToRemoteServer(t, terraformOptions, keyPair, username, "../../build/zarf-package-appliance-demo-doom.tar.zst", fmt.Sprintf("/home/%s/build/zarf-package-appliance-demo-doom.tar.zst", username), "0600")
	})

	teststructure.RunTestStage(t, "TEST", func() {
		terraformOptions := teststructure.LoadTerraformOptions(t, tmpFolder)
		keyPair := teststructure.LoadEc2KeyPair(t, tmpFolder)

		// Finally run the actual test
		testGameExample(t, terraformOptions, keyPair, username)
	})
}

func testGameExample(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair, username string) {
	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_instance_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: username,
	}

	// Make sure `zarf --help` doesn't error
	output, err := ssh.CheckSshCommandE(t, publicHost, fmt.Sprintf("sudo /home/%s/build/zarf --help", username))
	require.NoError(t, err, output)

	// run `zarf init`
	output, err = ssh.CheckSshCommandE(t, publicHost, fmt.Sprintf("sudo bash -c 'cd /home/%s/build && ./zarf init --confirm --components management --host 127.0.0.1'", username))
	require.NoError(t, err, output)

	// Wait until the Docker registry is ready
	output, err = ssh.CheckSshCommandE(t, publicHost, "timeout 300 bash -c 'while [[ \"$(curl -sfSL --retry 15 --retry-connrefused --retry-delay 5 -o /dev/null -w \"%{http_code}\" \"https://127.0.0.1/v2/\")\" != \"401\" ]]; do sleep 1; done' || false")
	require.NoError(t, err, output)

	// Deploy the game
	output, err = ssh.CheckSshCommandE(t, publicHost, fmt.Sprintf("sudo bash -c 'cd /home/%s/build && ./zarf package deploy zarf-package-appliance-demo-doom.tar.zst --confirm'", username))
	require.NoError(t, err, output)

	// Wait for the game to be live. Right now we're just checking that `curl` returns 0. It can be enhanced by scraping the HTML that gets returned or something.
	output, err = ssh.CheckSshCommandE(t, publicHost, "timeout 300 bash -c 'while [[ \"$(curl -sfSL --retry 15 --retry-connrefused --retry-delay 5 -o /dev/null -w \"%{http_code}\" \"https://127.0.0.1\")\" != \"200\" ]]; do sleep 1; done' || false")
	require.NoError(t, err, output)

	// Run `zarf destroy` to make sure that works correctly
	output, err = ssh.CheckSshCommandE(t, publicHost, fmt.Sprintf("sudo bash -c 'cd /home/%s/build && ./zarf destroy --confirm'", username))
	require.NoError(t, err, output)
}
