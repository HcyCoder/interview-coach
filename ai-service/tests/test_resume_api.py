import unittest

from httpx import ASGITransport, AsyncClient

from app.main import app


def build_pdf(text: str) -> bytes:
    """Build a minimal one-page PDF that contains the provided text.

    Inputs:
        text: Plain text to embed into the PDF content stream.
    Outputs:
        A valid PDF byte string that PyPDF2 can read and extract text from.
    """

    escaped = text.replace("\\", "\\\\").replace("(", "\\(").replace(")", "\\)")
    stream = f"BT /F1 12 Tf 72 720 Td ({escaped}) Tj ET".encode("latin-1")

    objects = [
        b"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n",
        b"2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n",
        (
            b"3 0 obj\n"
            b"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] "
            b"/Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>\n"
            b"endobj\n"
        ),
        b"4 0 obj\n<< /Length %d >>\nstream\n" % len(stream)
        + stream
        + b"\nendstream\nendobj\n",
        b"5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n",
    ]

    output = bytearray(b"%PDF-1.4\n")
    offsets = [0]
    for obj in objects:
        offsets.append(len(output))
        output.extend(obj)

    xref_offset = len(output)
    output.extend(f"xref\n0 {len(objects) + 1}\n".encode("latin-1"))
    output.extend(b"0000000000 65535 f \n")
    for offset in offsets[1:]:
        output.extend(f"{offset:010d} 00000 n \n".encode("latin-1"))
    output.extend(
        (
            "trailer\n"
            f"<< /Size {len(objects) + 1} /Root 1 0 R >>\n"
            f"startxref\n{xref_offset}\n%%EOF"
        ).encode("latin-1")
    )
    return bytes(output)


class ResumeApiTestCase(unittest.IsolatedAsyncioTestCase):
    async def asyncSetUp(self) -> None:
        self.client = AsyncClient(transport=ASGITransport(app=app), base_url="http://test")

    async def asyncTearDown(self) -> None:
        await self.client.aclose()

    async def test_parse_resume_returns_structured_profile(self) -> None:
        pdf = build_pdf(
            "Tech stack: Python, FastAPI, React, Go. "
            "Project: Interview Coach | Summary: built an AI interview practice platform."
        )

        response = await self.client.post(
            "/parse-resume",
            files={"file": ("resume.pdf", pdf, "application/pdf")},
        )

        self.assertEqual(response.status_code, 200)
        payload = response.json()
        self.assertEqual(payload["tech_stack"], ["Python", "FastAPI", "React", "Go"])
        self.assertEqual(payload["projects"][0]["name"], "Interview Coach")
        self.assertIn("AI interview practice platform", payload["summary"])

    async def test_generate_question_uses_resume_profile(self) -> None:
        response = await self.client.post(
            "/generate-question",
            json={
                "resume_profile": {
                    "tech_stack": ["Go", "MySQL", "Redis"],
                    "projects": [
                        {
                            "name": "Interview Coach",
                            "summary": "Built an AI interview practice platform",
                            "tech_stack": ["Go", "FastAPI"],
                            "highlights": ["shipped a candidate-facing flow"],
                        }
                    ],
                    "summary": "Strong backend and platform experience.",
                    "strengths": ["backend systems"],
                    "suggested_focus_areas": ["system design"],
                },
                "conversation": [
                    {"role": "assistant", "content": "Tell me about your background."},
                    {"role": "user", "content": "I led the backend for the resume parsing flow."},
                ],
            },
        )

        self.assertEqual(response.status_code, 200)
        payload = response.json()
        self.assertIn("trade-offs", payload["question"].lower())
        self.assertTrue(payload["question"].endswith("?"))


if __name__ == "__main__":
    unittest.main()
