using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class CameraCrl : MonoBehaviour
{
    public GameObject player;
    public float distanceOfPlayer = 10.0f;
    private Vector3 targetOffset;

    private float smoothTime = 0.01f;
    private Vector3 velocity = Vector3.zero;

    // Start is called before the first frame update
    void Start()
    {

    }

    // Update is called once per frame
    void Update()
    {
        transform.forward = player.transform.forward;
        float playerRotateY = player.transform.eulerAngles.y;
        float angleY = playerRotateY * Mathf.PI / 180.0f;
        targetOffset = new Vector3(-distanceOfPlayer * Mathf.Sin(angleY), 3, -distanceOfPlayer * Mathf.Cos(angleY));
        Vector3 targetPos = player.transform.position + targetOffset;

        transform.position = Vector3.SmoothDamp(transform.position, targetPos, ref velocity, smoothTime);
    }
}
