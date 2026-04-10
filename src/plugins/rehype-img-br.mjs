import { visit } from "unist-util-visit";

export function rehypeImgBr() {
	return (tree) => {
		visit(tree, "element", (node, index, parent) => {
			if (node.tagName === "img" && parent && index !== null) {
				const br = {
					type: "element",
					tagName: "br",
					properties: {},
					children: [],
				};
				parent.children.splice(index + 1, 0, br);
			}
		});
	};
}
