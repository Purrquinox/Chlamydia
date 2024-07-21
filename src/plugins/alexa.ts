import { PluginInterface, PluginType } from "../types.js";

class AlexaPlugin implements PluginInterface {
	private name: string;

	get(name: string) {
		this.name = name;
	}
}

const plugin: PluginType = {
	name: "Amazon Alexa",
	description: "Control your PC lighting using Amazon Alexa.",
	logo: "https://m.media-amazon.com/images/I/11e1mV2nO-L.png",
	show: false,
	plugin: new AlexaPlugin(),
};
export default plugin;
