from fastapi import APIRouter, File, HTTPException, UploadFile

from app.schemas import GenerateQuestionRequest, QuestionResponse, ResumeProfile
from app.services.question_service import generate_interview_question
from app.services.resume_service import build_resume_profile, extract_text_from_pdf

router = APIRouter()


@router.post("/parse-resume", response_model=ResumeProfile)
async def parse_resume(file: UploadFile = File(...)) -> ResumeProfile:
    """Parse an uploaded PDF resume into a structured JSON profile.

    Inputs:
        file: PDF resume uploaded as multipart form data.
    Outputs:
        Structured resume profile with tech stack, projects, and summary.
    """

    if file.content_type not in {"application/pdf", "application/x-pdf"}:
        raise HTTPException(status_code=415, detail="resume must be a PDF")

    pdf_bytes = await file.read()
    try:
        resume_text = extract_text_from_pdf(pdf_bytes)
    except ValueError as exc:
        raise HTTPException(status_code=400, detail=str(exc)) from exc

    return build_resume_profile(resume_text)


@router.post("/generate-question", response_model=QuestionResponse)
async def generate_question(payload: GenerateQuestionRequest) -> QuestionResponse:
    """Generate the next interview question from structured resume context.

    Inputs:
        payload: Structured resume profile and conversation history.
    Outputs:
        A single interview question.
    """

    question = generate_interview_question(payload.resume_profile, payload.conversation)
    return QuestionResponse(question=question)
