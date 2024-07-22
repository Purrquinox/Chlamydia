import { FastifyReply, FastifyRequest } from "fastify";
import { platforms } from "../../../main.js";
import { PlatformType } from "../../../types.js";

export default {
	url: "/get",
	method: "GET",
	schema: {
		summary: "Get Platforms",
		description: "This endpoint returns all platforms.",
		tags: ["platforms"],
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		let availablePlatforms: PlatformType[] = [];
		platforms.forEach((p) => availablePlatforms.push(p));

		return reply.status(200).send({
			platforms: availablePlatforms,
		});
	},
};
