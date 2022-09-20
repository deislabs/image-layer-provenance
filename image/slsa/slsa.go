/*
Copyright © 2022 Johnson Shi <Johnson.Shi@microsoft.com>
*/
package slsa

import (
	intoto "github.com/in-toto/in-toto-golang/in_toto"
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
)

func (s *ImageManifestLayerSlsaProvenance) GetImageManifestLayerSlsaProvenance() (*intoto.ProvenanceStatement, error) {
	return &intoto.ProvenanceStatement{
		StatementHeader: intoto.StatementHeader{
			Type:          intoto.StatementInTotoV01,
			PredicateType: slsa.PredicateSLSAProvenance,
			Subject: []intoto.Subject{
				{
					Name: string(s.LayerHistory.LayerDescriptor.Digest),
					Digest: slsa.DigestSet{
						string(s.LayerHistory.LayerDescriptor.Digest.Algorithm()): s.LayerHistory.LayerDescriptor.Digest.Encoded(),
					},
				},
			},
		},
		Predicate: slsa.ProvenancePredicate{
			Builder:   slsa.ProvenanceBuilder{ID: s.BuilderID},
			BuildType: s.BuildType,
			Invocation: slsa.ProvenanceInvocation{
				ConfigSource: slsa.ConfigSource{
					URI:        s.RepoURIContainingImageSource,
					Digest:     slsa.DigestSet{"commit": s.RepoGitCommit},
					EntryPoint: s.RepoPathToImageSource,
				},
				// NOTE: Invocation.Parameters is the only field that is customizable and allows for a custom JSON schema.
				Parameters: map[string]interface{}{
					"LayerHistory": s.LayerHistory,
				},
			},
			Metadata: &slsa.ProvenanceMetadata{
				BuildInvocationID: s.BuildInvocationID,
				BuildStartedOn:    s.BuildStartedOn,
				BuildFinishedOn:   s.BuildFinishedOn,
			},
		},
	}, nil
}
