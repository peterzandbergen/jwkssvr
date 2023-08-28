VERSION="10.0.2"

IMAGE=jwkssvr
REGISTRY="cir-cn.chp.belastingdienst.nl/zandp06"

REMOTE_TAGS=(
    "$REGISTRY/$IMAGE:$VERSION"
    "$REGISTRY/$IMAGE:latest"
)

LOCAL_TAGS=(
    "$IMAGE:$VERSION"
    "$IMAGE:latest"
)

# function get_tags() {
#     local tags=""
#     for t in ${REMOTE_TAGS[@]}
#     do
#         tags="$tags $(echo -n " --tag $t ")"
#     done
#     for t in ${ALL_TAGS[@]}
#     do
#         tags="$tags $(echo -n " --tag $t ")"
#     done
#     echo $tags
# }

function set_version() {
    VERSION="$1"

    REGISTRY="cir-cn.chp.belastingdienst.nl/zandp06"

    REMOTE_TAGS=(
        "$REGISTRY/$IMAGE:$VERSION"
        "$REGISTRY/$IMAGE:latest"
    )

    LOCAL_TAGS=(
        "$IMAGE:$VERSION"
        "$IMAGE:latest"
    )    
}

function get_tags() {
    local tags=""
    local arr=("$@")
    for t in "${arr[@]}"
    do
        tags="$tags $(echo -n " --tag $t ")"
    done
    echo -n $tags
}

function get_remote_tags() {
    get_tags "${REMOTE_TAGS[@]}"
}

function get_local_tags() {
    get_tags "${LOCAL_TAGS[@]}"
}

function get_all_tags() {
    echo -n $(get_remote_tags) $(get_local_tags)
}

function get_version_from_go_run() {
    (
        cd $WDIR/..
        echo -n $(go run ./cmd/jwkssvr -version)
    )
}

# REMOTE_IMAGE="cir-cn.chp.belastingdienst.nl/zandp06/jwkssvr:$VERSION"
# TAGS="--tag $REMOTE_IMAGE --tag jwkssvr:latest"
