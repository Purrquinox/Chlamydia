import { FastifyReply, FastifyRequest } from "fastify";
import platforms from "../../../platformIndex";
import {
	name,
	description,
	version,
	dependencies,
	devDependencies,
} from "../../../../package.json";
import { PlatformType } from "../../../types";

export default {
	url: "/",
	method: "GET",
	schema: {
		summary: "Hello",
		description:
			"This endpoint is the index page for our API, and lists some important information.",
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		let availablePlatforms: PlatformType[] = [];
		platforms.forEach((p) => availablePlatforms.push(p));

		return reply.status(200).send({
			name: name.charAt(0).toUpperCase() + name.slice(1),
			description,
			version,
			dependencies: { ...dependencies, ...devDependencies },
			documentation: "/docs",
			platforms: availablePlatforms,
			timestamp: new Date(),
		});
	},
};
