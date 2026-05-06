from app.schemas import ConversationTurn, ResumeProfile


def generate_interview_question(
    resume_profile: ResumeProfile, conversation: list[ConversationTurn]
) -> str:
    """Generate the next interview question from resume context and history.

    Inputs:
        resume_profile: Structured resume payload from the parser.
        conversation: Prior conversation turns for the session.
    Outputs:
        A single follow-up question suitable for streaming to the client.
    """

    project_name = (
        resume_profile.projects[0].name if resume_profile.projects else "your recent project"
    )
    primary_stack = ", ".join(resume_profile.tech_stack[:3]) or "your stack"
    last_user_message = next(
        (turn.content for turn in reversed(conversation) if turn.role == "user"),
        "",
    )

    if not last_user_message:
        return (
            f"Walk me through {project_name} and the biggest trade-offs you made "
            f"while working with {primary_stack}?"
        )

    return (
        f"You mentioned '{last_user_message}'. What trade-offs did you make, "
        f"and how would you improve the design if the traffic doubled?"
    )
