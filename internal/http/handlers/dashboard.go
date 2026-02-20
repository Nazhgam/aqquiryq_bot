package handlers

import (
	"fmt"
	"net/http"
)

func (h *handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	classes, err := h.contentRepo.GetAvailableClasses(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contents := make(map[string][]Presentation)

	for _, class := range classes {
		byClass, err := h.contentRepo.GetContentByClass(ctx, class)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		classContent := fmt.Sprintf("%d-сынып", class)

		presentationContents := make([]Presentation, len(byClass))
		for i, content := range byClass {
			presentationContents[i] = Presentation{
				ID:            fmt.Sprintf("%d", content.ID),
				Title:         content.Title,
				GroupName:     classContent,
				CanvaEmbedURL: content.CanvaURL,
			}
		}

		contents[classContent] = presentationContents
	}

	h.templates.ExecuteTemplate(w, "dashboard.html", contents)
}
