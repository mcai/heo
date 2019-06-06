package uncore

func  SetupCacheCoherenceFlowTree(flow CacheCoherenceFlow) {
	if flow.ProducerFlow() == nil {
		flow.SetAncestorFlow(flow)
		flow.Generator().MemoryHierarchy().SetPendingFlows(
			append(flow.Generator().MemoryHierarchy().PendingFlows(), flow),
		)
	} else {
		flow.SetAncestorFlow(flow.ProducerFlow().AncestorFlow())
		flow.ProducerFlow().SetChildFlows(append(flow.ProducerFlow().ChildFlows(), flow))
	}

	flow.AncestorFlow().SetNumPendingDescendantFlows(
		flow.AncestorFlow().NumPendingDescendantFlows() + 1)
}
