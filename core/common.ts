/**
 * This file contains common functions used throughout Chlamydia
 */

// Modules
import fs from "fs";
import path from "path";

/**
 * Function to get all files in a directory recursively
 */
export const getFilesInDirectory = (dir: string) => {
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
