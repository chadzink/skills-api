# Skill API Model Information
The purpose of this document is to describe the model used in this project and how they relate to eachother.

## Major Entities
The "major entities" in this project are:
* Skills: A skill that a person can have
* Person: A person who has skills
* Project: A proejct that a person applied a skill

## Supporting entities
There are supporting entities that provide context around the relationships of the major entities. These are:
* Categories: A classifiction level for skills to help organize
* Team: A grouping of persons during  aperiod of time
* Expertise: A scale that relates to a persons expierence in a skill

## Additional entities
TDB

## Entity Relationship Diagram (ERD)
The following ERD describes the current design of the database for the API
```mermaid
erDiagram
    categories {
        id bigint pk
        created_at timestamp
        updated_at timestamp
        deleted_at timestamp
        name text
        description text
        short_key text
        active boolean
    }

    skills {
        id bigint pk
        created_at timestamp
        updated_at timestamp
        deleted_at timestamp
        name text
        description text
        short_key text
        active boolean
    }

    people {
        id bigint pk
        created_at timestamp
        updated_at timestamp
        deleted_at timestamp
        email text
        name text
        phone text
        profile text
    }

    teams {
        id bigint pk
        created_at timestamp
        updated_at timestamp
        deleted_at timestamp
        name text
        active boolean
    }

    projects {
        id bigint pk
        created_at timestamp
        updated_at timestamp
        deleted_at timestamp
        name text
        start_dt datetime
        end_dt datetime
        team_id bigint fk
        active boolean
    }

    project_needs_expertise_in_skills {
        id bigint pk
        created_at timestamp
        updated_at timestamp
        deleted_at timestamp
        project_id  bigint fk
        expertise_id bigint fk
        skill_id bigint fk
    }

    teams ||--o{ projects : "team_id"
    projects ||--o{ project_needs_expertise_in_skills : "has many"
    project_needs_expertise_in_skills }o--|| expertises : "needs a"
    project_needs_expertise_in_skills }o--|| skills : "for a"

    expertises {
        id bigint pk
        created_at timestamp
        updated_at timestamp
        deleted_at timestamp
        name text
        description text
        order int
    }

    %% ----

    skill_category {
        category_id bigint fk
        skill_id bigint fk
    }

    categories ||--o{ skill_category : "category_id"
    skills ||--o{ skill_category : "skill_id"

    person_team {
        person_id bigint fk
        team_id bigint fk
        start_dt datetime
        end_dt datetime
    }

    people ||--o{ person_team : "person_id"
    teams ||--o{ person_team : "team_id"

    person_skills {
        person_id bigint fk
        skill_id bigint fk
        expertise_id bigint fk
        last_used datetime
    }

    people ||--o{ person_skills : "person_id"
    skills ||--o{ person_skills : "skill_id"
    expertises ||--o{ person_skills : "expertise_id"

    person_project {
        person_id bigint fk
        project_id bigint fk
        start_dt datetime
        end_dt datetime
    }

    people ||--o{ person_project : "person_id"
    projects ||--o{ person_project : "project_id"

```
