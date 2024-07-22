import { PluginInterface, PluginType } from "../types.js";

class TestPlugin implements PluginInterface {
	private name: string;

	get(name: string) {
		this.name = name;
	}
}

const plugin: PluginType = {
	name: "Test Plugin",
	description:
		"This is a test plugin, to ensure that everything works as expected without actually changing anything.",
	logo: "https://selectdev.purrquinox.com/logo.png",
	show: true,
	plugin: new TestPlugin(),
};
export default plugin;
