import { FastifyReply, FastifyRequest } from "fastify";
import { platforms, plugins } from "../../main.js";
import config from "../../config.js";
import { PlatformType, PluginType } from "../../types.js";

const { name, description, version, port } = config;

export default {
	url: "/",
	method: "GET",
	schema: {
		summary: "Hello",
		description:
			"This endpoint is the index page for our API, and lists some important information.",
		tags: [],
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		let supportedPlatforms: PlatformType[] = [];
		platforms.forEach((p) => {
			if (p.show) supportedPlatforms.push(p);
		});

		let supportedPlugins: PluginType[] = [];
		plugins.forEach((p) => {
			if (p.show) supportedPlugins.push(p);
		});

		return reply.status(200).send({
			name: name.charAt(0).toUpperCase() + name.slice(1),
			description,
			version,
			documentation: "/docs",
			supported_platforms: supportedPlatforms,
			supported_plugins: supportedPlugins,
			timestamp: new Date(),
		});
	},
};
