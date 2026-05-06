import re
from io import BytesIO

from PyPDF2 import PdfReader

from app.schemas import ResumeProfile, ResumeProject

TECH_KEYWORDS = [
    "Python",
    "FastAPI",
    "React",
    "Next.js",
    "TypeScript",
    "Go",
    "GORM",
    "MySQL",
    "Redis",
    "RabbitMQ",
    "Milvus",
    "LLM",
    "PyPDF2",
    "PostgreSQL",
    "OpenAI",
]


def extract_text_from_pdf(pdf_bytes: bytes) -> str:
    """Extract text from a PDF document.

    Inputs:
        pdf_bytes: Raw PDF bytes uploaded by the user.
    Outputs:
        Normalized extracted text with blank lines removed.
    Raises:
        ValueError: If the PDF cannot be parsed or contains no extractable text.
    """

    try:
        reader = PdfReader(BytesIO(pdf_bytes))
    except Exception as exc:  # noqa: BLE001
        raise ValueError("invalid pdf file") from exc

    text_parts: list[str] = []
    for page in reader.pages:
        extracted = page.extract_text() or ""
        cleaned = extracted.strip()
        if cleaned:
            text_parts.append(cleaned)

    text = "\n".join(text_parts).strip()
    if not text:
        raise ValueError("resume pdf contains no extractable text")

    return text


def build_resume_profile(resume_text: str) -> ResumeProfile:
    """Build a structured resume profile from extracted resume text.

    Inputs:
        resume_text: Plain text extracted from the candidate resume.
    Outputs:
        A structured profile suitable for downstream interview generation.
    """

    normalized_text = _normalize_text(resume_text)
    tech_stack = _extract_tech_stack(normalized_text)
    projects = _extract_projects(normalized_text, tech_stack)
    summary = _build_summary(normalized_text, tech_stack, projects)
    strengths = _build_strengths(tech_stack, projects)
    suggested_focus_areas = _build_focus_areas(tech_stack, projects)

    return ResumeProfile(
        tech_stack=tech_stack,
        projects=projects,
        summary=summary,
        strengths=strengths,
        suggested_focus_areas=suggested_focus_areas,
    )


def _normalize_text(text: str) -> str:
    return re.sub(r"\s+", " ", text).strip()


def _extract_tech_stack(text: str) -> list[str]:
    lower_text = text.lower()
    found: list[tuple[int, str]] = []
    for keyword in TECH_KEYWORDS:
        position = lower_text.find(keyword.lower())
        if position >= 0:
            found.append((position, keyword))

    found.sort(key=lambda item: item[0])
    ordered: list[str] = []
    seen: set[str] = set()
    for _, keyword in found:
        if keyword not in seen:
            seen.add(keyword)
            ordered.append(keyword)
    return ordered


def _extract_projects(text: str, tech_stack: list[str]) -> list[ResumeProject]:
    project_pattern = re.compile(
        r"(?:project|projects?)\s*[:\-]\s*(?P<name>[^|;\n]+)"
        r"(?:[|;]\s*(?P<rest>.*?))?(?=\s*(?:project|projects?)\s*[:\-]|$)",
        re.IGNORECASE,
    )
    matches = list(project_pattern.finditer(text))
    projects: list[ResumeProject] = []

    for match in matches:
        name = match.group("name").strip().rstrip(".")
        rest = match.group("rest") or ""
        summary = _capture_field(rest, ["summary", "description", "impact"])
        highlights = _capture_list_field(rest, ["highlights", "impact", "results"])
        project_tech_stack = _capture_list_field(rest, ["tech", "stack", "tech stack"])
        if not project_tech_stack:
            project_tech_stack = tech_stack[:3]
        if not summary:
            summary = f"Built {name} and shipped a candidate-facing experience."
        if not highlights:
            highlights = [summary]
        projects.append(
            ResumeProject(
                name=name,
                summary=summary,
                tech_stack=project_tech_stack,
                highlights=highlights,
            )
        )

    if projects:
        return projects

    fallback_name = "Resume Project"
    fallback_summary = "Built a product with practical engineering impact."
    if tech_stack:
        fallback_summary = f"Built systems with {', '.join(tech_stack[:3])}."

    return [
        ResumeProject(
            name=fallback_name,
            summary=fallback_summary,
            tech_stack=tech_stack[:3],
            highlights=[fallback_summary],
        )
    ]


def _capture_field(text: str, field_names: list[str]) -> str:
    for field_name in field_names:
        pattern = re.compile(
            rf"(?:{field_name})\s*:\s*(?P<value>[^|;]+)", re.IGNORECASE
        )
        match = pattern.search(text)
        if match:
            return match.group("value").strip().rstrip(".")
    return ""


def _capture_list_field(text: str, field_names: list[str]) -> list[str]:
    captured = _capture_field(text, field_names)
    if not captured:
        return []
    parts = [part.strip(" .") for part in re.split(r"[,/]", captured) if part.strip()]
    return [part for part in parts if part]


def _build_summary(
    text: str, tech_stack: list[str], projects: list[ResumeProject]
) -> str:
    if projects:
        project_names = ", ".join(project.name for project in projects[:2])
        project_context = projects[0].summary
    else:
        project_names = "candidate projects"
        project_context = "broad product and engineering delivery"

    primary_stack = ", ".join(tech_stack[:4]) if tech_stack else "multiple systems"
    return (
        f"Candidate profile centers on {project_names}; {project_context} while "
        f"working across {primary_stack}."
    )


def _build_strengths(tech_stack: list[str], projects: list[ResumeProject]) -> list[str]:
    strengths: list[str] = []
    if "Go" in tech_stack or "MySQL" in tech_stack:
        strengths.append("backend systems")
    if "FastAPI" in tech_stack or "Python" in tech_stack:
        strengths.append("service integration")
    if "React" in tech_stack or "Next.js" in tech_stack:
        strengths.append("frontend delivery")
    if projects:
        strengths.append("project communication")
    return strengths or ["general engineering"]


def _build_focus_areas(
    tech_stack: list[str], projects: list[ResumeProject]
) -> list[str]:
    areas: list[str] = []
    if "Go" in tech_stack or "MySQL" in tech_stack:
        areas.append("system design")
    if "FastAPI" in tech_stack or "Python" in tech_stack:
        areas.append("API design")
    if "React" in tech_stack or "Next.js" in tech_stack:
        areas.append("frontend trade-offs")
    if projects:
        areas.append("project depth")
    return areas or ["technical fundamentals"]
