// routers/quiz.go
package routers

import (
    "bytes"
    "html/template"
    "log"

    "github.com/gofiber/fiber/v2"
)

type Question struct {
    Text    string
    Options []string
}

type QuizData struct {
    Title     string
    Questions []Question
}

func SetupQuizRoutes(app *fiber.App) {
    // Parse only header.html and quizForm.html, using QuizForm as the root
    tmpl := template.Must(
        template.New("quizForm").
            ParseFiles(
                "./views/header.html",   // registers "header.html"
                "./views/quizForm.html", // registers "quizForm.html"
            ),
    )

    app.Get("/quiz", func(c *fiber.Ctx) error {
        quiz := QuizData{
            Title: "Trivia Quiz",
            Questions: []Question{
                {"What is the capital of France?", []string{"Paris", "Berlin", "Madrid", "Rome"}},
                {"What is 2 + 2?", []string{"3", "4", "5", "6"}},
            },
        }

        var buf bytes.Buffer
        // Execute the template defined in quizForm.html (which uses {{define "quizForm"}}…)
        if err := tmpl.ExecuteTemplate(&buf, "quizForm", quiz); err != nil {
            log.Println("Template exec error:", err)
            return c.Status(500).SendString("Template execution error: " + err.Error())
        }
        return c.Type("html").Send(buf.Bytes())
    })

    app.Post("/quiz", func(c *fiber.Ctx) error {
        answers := make(map[string]string)
        c.Request().PostArgs().VisitAll(func(k, v []byte) {
            answers[string(k)] = string(v)
        })
        return c.JSON(fiber.Map{"status": "received", "answers": answers})
    })
}
