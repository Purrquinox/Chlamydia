// Packages
import { PlatformType, PluginType } from "./types.js";
import { getFilesInDirectory } from "./common.js";

// Maps
const platforms: Map<string, PlatformType> = new Map();
const plugins: Map<string, PluginType> = new Map();

// Get the data from their folders.
const platformFiles = getFilesInDirectory("./dist/platforms").filter((file) =>
	file.endsWith(".js")
);
const pluginFiles = getFilesInDirectory("./dist/plugins").filter((file) =>
	file.endsWith(".js")
);

// Loop thru the files, and import them to receive data.
for (const file of platformFiles) {
	import(`../${file}`)
		.then(async (i) => {
			const module = i.default;
			platforms.set(module.name, module);
		})
		.catch((error) => {
			console.error(`Error importing ${file}: ${error}`);
		});
}

for (const file of pluginFiles) {
	import(`../${file}`)
		.then(async (i) => {
			const plugin = i.default;
			plugins.set(plugin.name, plugin);
		})
		.catch((error) => {
			console.error(`Error importing ${file}: ${error}`);
		});
}

// Export
export { platforms, plugins };
