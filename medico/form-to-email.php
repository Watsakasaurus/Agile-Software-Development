
<?php
/*
if(!isset($_POST['submit']))
{
	//This page should not be accessed directly. Need to submit the form.
	echo "error; you need to submit the form!";
}
$name = $_POST['name'];
$visitor_email = $_POST['email'];
$message = $_POST['message'];

//Validate first
if(empty($name)||empty($visitor_email)) 
{
    echo "Name and email are mandatory!";
    exit;
}

if(IsInjected($visitor_email))
{
    echo "Bad email value!";
    exit;
}

$email_from = 'amy.gourlay@ntlworld.com';//<== update the email address
$email_subject = "New Form submission";
$email_body = "You have received a new message from the user $name.\n".
    "Here is the message:\n $message".
    
$to = "amy.gourlay@ntlworld.com";//<== update the email address
$headers = "From: $email_from \r\n";
$headers .= "Reply-To: $visitor_email \r\n";
//Send the email!
mail($to,$email_subject,$email_body,$headers);
//done. redirect to thank-you page.
header('Location: thank-you.html');


// Function to validate against any email injection attempts
function IsInjected($str)
{
  $injections = array('(\n+)',
              '(\r+)',
              '(\t+)',
              '(%0A+)',
              '(%0D+)',
              '(%08+)',
              '(%09+)'
              );
  $inject = join('|', $injections);
  $inject = "/$inject/i";
  if(preg_match($inject,$str))
    {
    return true;
  }
  else
    {
    return false;
  }
}
   */

  
  error_reporting(-1);
  
  if(isset($_POST['submit']))
  {
  $name = $_POST['name']; 
  $submit_links = $_POST['submit_links']; 
  $from_add = "amy.gourlay@ntlworld.com"; 
  $to_add = "amy.gourlay@ntlworld.com"; 
  $subject = "Your Subject Name";
  $message = "Name:$name \n Sites: $submit_links";
  
  $headers = 'From: submit@webdesignrepo.com' . "\r\n" .
  
  'Reply-To: ben@webdesignrepo.com' . "\r\n";
  
  if(mail($to_add,$subject,$message,$headers)) 
  {
      $msg = "Mail sent";
  
   echo $msg;
  } 
  
  print "<p>Thanks $name</p>" ;
  }
  
  // else conditional statement for if(isset($_POST['submit']))
  else {
  echo "Sorry, you cannot do that from here. Please fill in the form first.";
  }
  





?> 