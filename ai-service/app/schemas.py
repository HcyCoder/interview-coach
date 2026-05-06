from typing import Literal

from pydantic import BaseModel, Field


class ResumeProject(BaseModel):
    """ResumeProject describes one concrete project in the candidate profile."""

    name: str
    summary: str
    tech_stack: list[str] = Field(default_factory=list)
    highlights: list[str] = Field(default_factory=list)


class ResumeProfile(BaseModel):
    """ResumeProfile is the structured resume payload returned by the parser."""

    tech_stack: list[str] = Field(default_factory=list)
    projects: list[ResumeProject] = Field(default_factory=list)
    summary: str = ""
    strengths: list[str] = Field(default_factory=list)
    suggested_focus_areas: list[str] = Field(default_factory=list)


class ConversationTurn(BaseModel):
    """ConversationTurn stores one message in the interview history."""

    role: Literal["assistant", "user"]
    content: str


class GenerateQuestionRequest(BaseModel):
    """GenerateQuestionRequest asks the AI service for the next interview question."""

    resume_profile: ResumeProfile
    conversation: list[ConversationTurn] = Field(default_factory=list)


class QuestionResponse(BaseModel):
    """QuestionResponse returns the next interview question."""

    question: str
