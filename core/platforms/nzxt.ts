import { PlatformInterface, PlatformType } from "../types.js";

class NZXT implements PlatformInterface {
	private name: string;

	get(name: string) {
		this.name = name;
	}
}

const platform: PlatformType = {
	name: "NZXT",
	description:
		"NZXT is a leading force in PC hardware, known for its groundbreaking designs and high-performance components. Their products, ranging from sleek cases to advanced cooling solutions, are crafted to elevate both the aesthetics and functionality of your gaming and workstation setups.",
	logo: "https://nzxt.com/assets/cms/34299/1611033291-nzxt-logo.png?auto=format&fit=max&h=540&w=540",
	show: false,
	platform: new NZXT(),
};
export default platform;
