import { FastifyReply, FastifyRequest } from "fastify";
import { plugins } from "../../../main.js";
import { PluginType } from "../../../types.js";

export default {
	url: "/get",
	method: "GET",
	schema: {
		summary: "Get Plugins",
		description: "This endpoint returns all plugins.",
		tags: ["plugins"],
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		let availablePlugins: PluginType[] = [];
		plugins.forEach((p) => availablePlugins.push(p));

		return reply.status(200).send({
			plugins: availablePlugins,
		});
	},
};
