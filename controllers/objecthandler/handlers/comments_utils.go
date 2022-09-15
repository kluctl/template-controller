package handlers

import (
	"github.com/kluctl/template-controller/controllers"
	"github.com/kluctl/template-controller/controllers/objecthandler/webgit"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func reconcileComment(clusterId string, mr webgit.MergeRequestInterface, tag string, obj client.Object, comment string, noteId *string, lastPostedBodyHash *string) error {
	var err error
	body := generateMarkerComment(tag, clusterId, obj.GetNamespace(), obj.GetName()) + "\n" + comment

	var existingNote webgit.Note
	if *noteId == "" {
		existingNote, err = findNote(clusterId, mr, tag, obj)
		if err != nil {
			return err
		}
		if existingNote != nil {
			*noteId = existingNote.GetId()
			*lastPostedBodyHash = ""
		} else {
			existingNote, err = mr.CreateMergeRequestNote(body)
			if err != nil {
				return err
			}
			*noteId = existingNote.GetId()
			*lastPostedBodyHash = controllers.Sha256String(body)
		}
	} else {
		var resp *gitlab.Response
		existingNote, err = mr.GetMergeRequestNote(*noteId)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				*noteId = ""
				*lastPostedBodyHash = ""
			}
			return err
		}
	}

	if *lastPostedBodyHash == controllers.Sha256String(body) {
		return nil
	}

	err = existingNote.UpdateBody(body)
	if err != nil {
		*noteId = ""
		*lastPostedBodyHash = ""
		return err
	}
	*lastPostedBodyHash = controllers.Sha256String(body)
	return nil
}

func findNote(clusterId string, mr webgit.MergeRequestInterface, tag string, obj client.Object) (webgit.Note, error) {
	notes, err := mr.ListMergeRequestNotes()
	if err != nil {
		return nil, err
	}
	for _, n := range notes {
		if !hasMarkerComment(n.GetBody(), tag, clusterId, obj.GetNamespace(), obj.GetName()) {
			continue
		}
		return n, nil
	}
	return nil, nil
}
