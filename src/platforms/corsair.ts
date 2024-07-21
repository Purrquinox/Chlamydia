import { PlatformInterface, PlatformType } from "../types.js";

class Corsair implements PlatformInterface {
	private name: string;

	get(name: string) {
		this.name = name;
	}
}

const platform: PlatformType = {
	name: "Corsair",
	description:
		"CORSAIR is a leading global developer and manufacturer of high-performance gear and technology for gamers, content creators, and PC enthusiasts.",
	logo: "https://res.cloudinary.com/corsair-pwa/image/upload/f_auto,q_auto/v1665096094/akamai/content/images/reusable/CORSAIR_logo_horizontal_white.png",
	show: false,
	platform: new Corsair(),
};
export default platform;
