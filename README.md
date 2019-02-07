# db-container

This Repo contains the main "database" running an http-api exposing the quadtree.

## API-Endpoints

| Endpoint | Description | POST parameters |
| --- | --- | --- |
| `"/"` | Index | |
| `"/new"` | Create a new star | `w` |
| `"/printall"` | Print all the trees in json| |
| `"/insert/{treeindex}"` | Insert the given star into the selected tree | `x`, `y`, `vx`, `vy`, `m` |
| `"/starlist/{treeindex}"` | List all the stars in the selected tree| |
| `"/dumptree/{treeindex}"` | Dump the json of the selected tree | |
| `"/updatetotalmass/{treeindex}"` | Update the total mass in the selected tree | |
| `"/updatecenterofmass/{treeindex}"` | Update the center of mass in the selected tree | |
| `"/metrics"` | Get the overall metrics | |
| `"/export/{treeindex}"` | Export the selected tree to `db/{treeindex}.json` | |
| `"/fastinsert/{filename}"` | Insert the selected file into a new tree | |

## Tables

### nodes

```postgresql
-- Table: public.nodes

-- DROP TABLE public.nodes;

CREATE TABLE public.nodes
(
    node_id bigint NOT NULL DEFAULT nextval('nodes_node_id_seq'::regclass),
    box_center_x numeric,
    box_center_y numeric,
    box_width numeric,
    center_of_mass_x numeric,
    center_of_mass_y numeric,
    total_mass numeric,
    depth numeric,
    star_id bigint,
    "subnode_A" bigint,
    "subnode_B" bigint,
    "subnode_C" bigint,
    "subnode_D" bigint,
    CONSTRAINT nodes_pkey PRIMARY KEY (node_id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.nodes
    OWNER to postgres;
```


