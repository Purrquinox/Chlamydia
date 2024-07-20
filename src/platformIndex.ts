/* Packages */
import fs from "node:fs";
import { PlatformType } from "./types";
import path from "path";

// Platform Map
const platforms: Map<string, PlatformType> = new Map();

// Function that gets all files from a directory
const getFilesInDirectory = (dir: string) => {
	let files: string[] = [];
	const filesInDir = fs.readdirSync(dir);

	for (const file of filesInDir) {
		const filePath = path.join(dir, file);
		const stat = fs.statSync(filePath);

		if (stat.isDirectory())
			files = files.concat(getFilesInDirectory(filePath));
		else files.push(filePath);
	}

	return files;
};

// Get files from the dist/platforms directory, and filter it to only get the .js files as needed.
const platformFiles  = getFilesInDirectory("./dist/platforms").filter(
	(file) => file.endsWith(".js")
);

// Loop thru the files, and import them to receive data.
for (const file of platformFiles) {
	import(`../${file}`)
		.then(async (i) => {
            const module = i.default.default;
            if (module.data.show) platforms.set(i.name, i);
		})
		.catch((error) => {
			console.error(`Error importing ${file}: ${error}`);
		});
}