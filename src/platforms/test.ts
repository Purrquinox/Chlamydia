import { PlatformInterface, PlatformType } from "../types";

class TestPlatform implements PlatformInterface {
	private name: string;

	get(name: string) {
		this.name = name;
	}
}

const platform: PlatformType = {
	name: "Test Platform",
	description:
		"This is a test platform, to ensure that everything works as expected without actually changing anything.",
	logo: "https://selectdev.purrquinox.com/logo.png",
	show: true,
	platform: new TestPlatform(),
};
export default platform;
